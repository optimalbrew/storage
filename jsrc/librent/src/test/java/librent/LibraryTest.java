/*
 * This Java source file was generated by the Gradle 'init' task.
 */
package librent;

import org.junit.Test;
import static org.junit.Assert.*;

import java.math.BigInteger;

public class LibraryTest {
    
    @Test public void testSomeLibraryMethod() {
        Library classUnderTest = new Library();
        assertFalse("someLibraryMethod should return 'true'", classUnderTest.someLibraryMethod());
    }
    
    @Test public void testAppHasAGreeting() {
        Library classUnderTest = new Library();
        assertNotNull("Lib should have a greeting", classUnderTest.getGreeting());
    } 

    @Test public void testReturnRent(){
        Library classUnderTest = new Library();
        BigInteger s = BigInteger.valueOf(4).pow(5);
        assertEquals("Rent is not equal to constant", s, classUnderTest.returnRent());
    }

}
